import { FormEvent, useEffect, useMemo, useState } from "react";

type Role = "tenant" | "landlord";
type View = "catalog" | "favorites" | "create" | "chats" | "profile";

type User = {
  id: number;
  name: string;
  email: string;
  phone: string;
  role: Role;
  avatarUrl?: string;
  bio?: string;
};

type Apartment = {
  id: number;
  title: string;
  description: string;
  type: string;
  price: number;
  district: string;
  address: string;
  rooms: number;
  floor: number;
  hasFurniture: boolean;
  hasWifi: boolean;
  hasWasher: boolean;
  ownerId: number;
  photos: string[];
};

type Conversation = {
  id: number;
  apartment?: Apartment;
  tenant?: User;
  landlord?: User;
};

type Message = {
  id: number;
  text: string;
  senderId: number;
  createdAt?: string;
};

type Filters = {
  q: string;
  district: string;
  minPrice: string;
  maxPrice: string;
  rooms: string;
  hasWifi: boolean;
  hasFurniture: boolean;
  hasWasher: boolean;
  sort: string;
};

const API_BASE = import.meta.env.VITE_API_URL || "/api";

const fallbackApartments: Apartment[] = [
  {
    id: 101,
    title: "2-комнатная квартира у Ботанического сада",
    description: "Светлая квартира с мебелью, Wi-Fi и удобным доступом к университетам.",
    type: "apartment",
    price: 260000,
    district: "Есиль",
    address: "Кабанбай Батыра 48",
    rooms: 2,
    floor: 8,
    hasFurniture: true,
    hasWifi: true,
    hasWasher: true,
    ownerId: 1,
    photos: [
      "https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?auto=format&fit=crop&w=1200&q=80",
    ],
  },
  {
    id: 102,
    title: "Студия в районе Expo",
    description: "Минималистичная студия для студента или молодого специалиста.",
    type: "studio",
    price: 190000,
    district: "Нура",
    address: "Туран 55",
    rooms: 1,
    floor: 12,
    hasFurniture: true,
    hasWifi: true,
    hasWasher: false,
    ownerId: 2,
    photos: [
      "https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?auto=format&fit=crop&w=1200&q=80",
    ],
  },
  {
    id: 103,
    title: "Комната рядом с университетом",
    description: "Аккуратная комната в спокойном районе, все базовое уже есть.",
    type: "room",
    price: 95000,
    district: "Алматы",
    address: "Абылайхана 12",
    rooms: 1,
    floor: 3,
    hasFurniture: true,
    hasWifi: true,
    hasWasher: true,
    ownerId: 3,
    photos: [
      "https://images.unsplash.com/photo-1493809842364-78817add7ffb?auto=format&fit=crop&w=1200&q=80",
    ],
  },
];

const districts = ["Есиль", "Алматы", "Сарыарка", "Байконур", "Нура", "Сарайшык"];

const initialFilters: Filters = {
  q: "",
  district: "",
  minPrice: "",
  maxPrice: "",
  rooms: "",
  hasWifi: false,
  hasFurniture: false,
  hasWasher: false,
  sort: "newest",
};

function App() {
  const [view, setView] = useState<View>("catalog");
  const [apartments, setApartments] = useState<Apartment[]>(fallbackApartments);
  const [favorites, setFavorites] = useState<Apartment[]>([]);
  const [filters, setFilters] = useState<Filters>(initialFilters);
  const [selected, setSelected] = useState<Apartment | null>(null);
  const [user, setUser] = useState<User | null>(() => readStoredUser());
  const [token, setToken] = useState(() => localStorage.getItem("alarent_access") || "");
  const [notice, setNotice] = useState("Готово к поиску жилья в Астане");
  const [authMode, setAuthMode] = useState<"login" | "register">("login");
  const [loading, setLoading] = useState(false);
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeConversation, setActiveConversation] = useState<number | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);

  useEffect(() => {
    void loadApartments();
  }, []);

  useEffect(() => {
    if (token) {
      void loadMe();
      void loadFavorites();
      void loadConversations();
    }
  }, [token]);

  const filteredFallback = useMemo(() => {
    if (apartments !== fallbackApartments) return apartments;
    return apartments.filter((item) => {
      const text = `${item.title} ${item.description} ${item.address}`.toLowerCase();
      if (filters.q && !text.includes(filters.q.toLowerCase())) return false;
      if (filters.district && item.district !== filters.district) return false;
      if (filters.minPrice && item.price < Number(filters.minPrice)) return false;
      if (filters.maxPrice && item.price > Number(filters.maxPrice)) return false;
      if (filters.rooms && item.rooms !== Number(filters.rooms)) return false;
      if (filters.hasWifi && !item.hasWifi) return false;
      if (filters.hasFurniture && !item.hasFurniture) return false;
      if (filters.hasWasher && !item.hasWasher) return false;
      return true;
    });
  }, [apartments, filters]);

  async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers = new Headers(options.headers);
    if (!(options.body instanceof FormData)) {
      headers.set("Content-Type", "application/json");
    }
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const response = await fetch(`${API_BASE}${path}`, { ...options, headers });
    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
      throw new Error(data.error || "Запрос не выполнен");
    }
    return data as T;
  }

  async function loadApartments() {
    setLoading(true);
    try {
      const params = buildFilterParams(filters);
      const data = await request<unknown[]>(`/apartaments${params}`);
      setApartments(data.map(normalizeApartment));
      setNotice("Каталог обновлен");
    } catch {
      setApartments(fallbackApartments);
      setNotice("Backend не отвечает, показываю демо-карточки");
    } finally {
      setLoading(false);
    }
  }

  async function loadMe() {
    try {
      const data = await request<Record<string, unknown>>("/me");
      const nextUser = normalizeUser(data);
      setUser(nextUser);
      localStorage.setItem("alarent_user", JSON.stringify(nextUser));
    } catch {
      logoutLocal();
    }
  }

  async function loadFavorites() {
    try {
      const data = await request<unknown[]>("/me/favorites");
      setFavorites(
        data
          .map((item) => normalizeApartment(read(item, "Apartment", "apartment") || item))
          .filter((apartment) => apartment.id > 0 && apartment.price > 0),
      );
    } catch {
      setFavorites([]);
    }
  }

  async function loadConversations() {
    try {
      const data = await request<unknown[]>("/conversations");
      setConversations(data.map(normalizeConversation));
    } catch {
      setConversations([]);
    }
  }

  async function loadMessages(conversationId: number) {
    setActiveConversation(conversationId);
    try {
      const data = await request<unknown[]>(`/conversations/${conversationId}/messages`);
      setMessages(data.map(normalizeMessage));
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function handleAuth(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const payload = Object.fromEntries(form.entries());
    const path = authMode === "login" ? "/auth/login" : "/auth/register";

    try {
      if (authMode === "register") {
        await request(path, {
          method: "POST",
          body: JSON.stringify(payload),
        });
        setAuthMode("login");
        setNotice("Аккаунт создан, теперь войди");
        return;
      }

      const data = await request<Record<string, unknown>>(path, {
        method: "POST",
        body: JSON.stringify(payload),
      });
      const access = String(read(data, "access_token", "token") || "");
      const refresh = String(read(data, "refresh_token") || "");
      const nextUser = normalizeUser(read(data, "user") || {});
      localStorage.setItem("alarent_access", access);
      localStorage.setItem("alarent_refresh", refresh);
      localStorage.setItem("alarent_user", JSON.stringify(nextUser));
      setToken(access);
      setUser(nextUser);
      setNotice(`Добро пожаловать, ${nextUser.name}`);
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function createApartment(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formElement = event.currentTarget;
    const form = new FormData(formElement);
    const photoUrls = String(form.get("photo_urls") || "")
      .split("\n")
      .map((url) => url.trim())
      .filter(Boolean);

    const payload = {
      title: form.get("title"),
      description: form.get("description"),
      type: form.get("type"),
      price: Number(form.get("price")),
      district: form.get("district"),
      address: form.get("address"),
      rooms: Number(form.get("rooms")),
      floor: Number(form.get("floor")),
      has_furniture: form.get("has_furniture") === "on",
      has_wifi: form.get("has_wifi") === "on",
      has_washer: form.get("has_washer") === "on",
      photo_urls: photoUrls,
    };

    try {
      await request("/apartaments", { method: "POST", body: JSON.stringify(payload) });
      setNotice("Объявление опубликовано");
      formElement.reset();
      setView("catalog");
      await loadApartments();
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function updateProfile(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const payload = {
      name: form.get("name"),
      phone: form.get("phone"),
      avatar_url: form.get("avatar_url"),
      bio: form.get("bio"),
    };
    try {
      const data = await request<Record<string, unknown>>("/me", {
        method: "PATCH",
        body: JSON.stringify(payload),
      });
      const nextUser = normalizeUser(data);
      setUser(nextUser);
      localStorage.setItem("alarent_user", JSON.stringify(nextUser));
      setNotice("Профиль обновлен");
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function toggleFavorite(apartment: Apartment) {
    if (!token) {
      setNotice("Сначала войди в аккаунт");
      return;
    }

    const isFavorite = favorites.some((item) => item.id === apartment.id);
    try {
      await request(`/apartaments/${apartment.id}/favorite`, {
        method: isFavorite ? "DELETE" : "POST",
      });
      await loadFavorites();
      setNotice(isFavorite ? "Удалено из избранного" : "Добавлено в избранное");
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function openConversation(apartment: Apartment) {
    if (!token) {
      setNotice("Сначала войди в аккаунт");
      return;
    }
    try {
      const data = await request<Record<string, unknown>>(`/apartaments/${apartment.id}/conversation`, {
        method: "POST",
      });
      const conversation = normalizeConversation(data);
      await loadConversations();
      setView("chats");
      await loadMessages(conversation.id);
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function deleteApartment(apartment: Apartment) {
    if (!token) {
      setNotice("Сначала войди в аккаунт");
      return;
    }
    if (!window.confirm("Удалить это объявление?")) {
      return;
    }

    try {
      await request(`/apartaments/${apartment.id}`, { method: "DELETE" });
      setSelected(null);
      setFavorites((items) => items.filter((item) => item.id !== apartment.id));
      setApartments((items) => items.filter((item) => item.id !== apartment.id));
      setNotice("Объявление удалено");
      await loadApartments();
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  async function sendMessage(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!activeConversation) return;
    const formElement = event.currentTarget;
    const form = new FormData(formElement);
    const text = String(form.get("text") || "").trim();
    if (!text) return;

    try {
      await request(`/conversations/${activeConversation}/messages`, {
        method: "POST",
        body: JSON.stringify({ text }),
      });
      formElement.reset();
      await loadMessages(activeConversation);
    } catch (error) {
      setNotice(getErrorMessage(error));
    }
  }

  function logoutLocal() {
    localStorage.removeItem("alarent_access");
    localStorage.removeItem("alarent_refresh");
    localStorage.removeItem("alarent_user");
    setToken("");
    setUser(null);
    setFavorites([]);
    setConversations([]);
    setMessages([]);
    setNotice("Ты вышел из аккаунта");
  }

  const listing = filteredFallback;

  if (view === "profile") {
    return (
      <main>
        <Header
          user={user}
          view={view}
          setView={setView}
          logout={logoutLocal}
        />
        <Profile
          user={user}
          authMode={authMode}
          setAuthMode={(mode) => {
            setAuthMode(mode);
            setNotice("");
          }}
          handleAuth={handleAuth}
          updateProfile={updateProfile}
          notice={notice}
        />
      </main>
    );
  }

  return (
    <main>
      <Header
        user={user}
        view={view}
        setView={setView}
        logout={logoutLocal}
      />

      <section className="hero">
        <div className="heroMedia" aria-hidden="true" />
        <div className="heroContent">
          <p className="eyebrow">AlaRent · аренда жилья в Астане</p>
          <h1>Найди квартиру без шума, лишних звонков и потерянных вариантов.</h1>
          <p className="heroText">
            Минималистичный каталог для студентов и молодых специалистов: поиск, фильтры,
            избранное и быстрый чат с арендодателем.
          </p>
          <div className="heroSearch">
            <input
              value={filters.q}
              onChange={(event) => setFilters({ ...filters, q: event.target.value })}
              placeholder="Район, ЖК, улица или описание"
            />
            <button onClick={loadApartments}>{loading ? "Ищу..." : "Найти"}</button>
          </div>
        </div>
      </section>

      <div className="shell">
        <StatusBar notice={notice} count={listing.length} />

        {view === "catalog" && (
          <Catalog
            filters={filters}
            setFilters={setFilters}
            apartments={listing}
            favorites={favorites}
            refresh={loadApartments}
            select={setSelected}
            toggleFavorite={toggleFavorite}
          />
        )}

        {view === "favorites" && (
          <ApartmentGrid
            title="Избранное"
            subtitle="Все сохраненные варианты в одном месте."
            apartments={favorites}
            favorites={favorites}
            select={setSelected}
            toggleFavorite={toggleFavorite}
          />
        )}

        {view === "create" && (
          <CreateApartmentForm user={user} onSubmit={createApartment} />
        )}

        {view === "chats" && (
          <Chats
            conversations={conversations}
            activeConversation={activeConversation}
            messages={messages}
            openConversation={loadMessages}
            sendMessage={sendMessage}
            user={user}
          />
        )}

      </div>

      {selected && (
        <ApartmentDetails
          apartment={selected}
          isFavorite={favorites.some((item) => item.id === selected.id)}
          close={() => setSelected(null)}
          toggleFavorite={toggleFavorite}
          openConversation={openConversation}
          deleteApartment={deleteApartment}
          user={user}
        />
      )}
    </main>
  );
}

function Header({ user, view, setView, logout }: {
  user: User | null;
  view: View;
  setView: (view: View) => void;
  logout: () => void;
}) {
  const links: { label: string; view: View }[] = [
    { label: "Каталог", view: "catalog" },
    { label: "Избранное", view: "favorites" },
    { label: "Подать", view: "create" },
    { label: "Чаты", view: "chats" },
    { label: "Профиль", view: "profile" },
  ];

  return (
    <header className="topbar">
      <button className="brand" onClick={() => setView("catalog")}>
        <span>A</span>
        AlaRent
      </button>
      <nav>
        {links.map((link) => (
          <button
            key={link.view}
            className={view === link.view ? "active" : ""}
            onClick={() => setView(link.view)}
          >
            {link.label}
          </button>
        ))}
      </nav>
      <div className="userPill">
        {user ? (
          <>
            <span>{user.name}</span>
            <button onClick={logout}>Выйти</button>
          </>
        ) : (
          <button onClick={() => setView("profile")}>Войти</button>
        )}
      </div>
    </header>
  );
}

function StatusBar({ notice, count }: { notice: string; count: number }) {
  return (
    <section className="statusBar">
      <span>{notice}</span>
      <strong>{count} вариантов</strong>
    </section>
  );
}

function Catalog(props: {
  filters: Filters;
  setFilters: (filters: Filters) => void;
  apartments: Apartment[];
  favorites: Apartment[];
  refresh: () => void;
  select: (apartment: Apartment) => void;
  toggleFavorite: (apartment: Apartment) => void;
}) {
  return (
    <section className="layout">
      <aside className="filters">
        <h2>Фильтры</h2>
        <label>
          Район
          <select
            value={props.filters.district}
            onChange={(event) => props.setFilters({ ...props.filters, district: event.target.value })}
          >
            <option value="">Все районы</option>
            {districts.map((district) => (
              <option key={district} value={district}>{district}</option>
            ))}
          </select>
        </label>
        <div className="twoColumns">
          <label>
            Цена от
            <input value={props.filters.minPrice} onChange={(event) => props.setFilters({ ...props.filters, minPrice: event.target.value })} />
          </label>
          <label>
            до
            <input value={props.filters.maxPrice} onChange={(event) => props.setFilters({ ...props.filters, maxPrice: event.target.value })} />
          </label>
        </div>
        <label>
          Комнаты
          <select value={props.filters.rooms} onChange={(event) => props.setFilters({ ...props.filters, rooms: event.target.value })}>
            <option value="">Любое</option>
            <option value="1">1</option>
            <option value="2">2</option>
            <option value="3">3</option>
            <option value="4">4+</option>
          </select>
        </label>
        <label>
          Сортировка
          <select value={props.filters.sort} onChange={(event) => props.setFilters({ ...props.filters, sort: event.target.value })}>
            <option value="newest">Сначала новые</option>
            <option value="price_asc">Дешевле</option>
            <option value="price_desc">Дороже</option>
            <option value="oldest">Сначала старые</option>
          </select>
        </label>
        <div className="checkStack">
          <label><input type="checkbox" checked={props.filters.hasFurniture} onChange={(event) => props.setFilters({ ...props.filters, hasFurniture: event.target.checked })} /> Мебель</label>
          <label><input type="checkbox" checked={props.filters.hasWifi} onChange={(event) => props.setFilters({ ...props.filters, hasWifi: event.target.checked })} /> Wi-Fi</label>
          <label><input type="checkbox" checked={props.filters.hasWasher} onChange={(event) => props.setFilters({ ...props.filters, hasWasher: event.target.checked })} /> Стиралка</label>
        </div>
        <button className="wide" onClick={props.refresh}>Показать результаты</button>
      </aside>

      <ApartmentGrid
        title="Аренда жилья"
        subtitle="Карточки в стиле маркетплейса: фото, цена, район и важные удобства сразу видны."
        apartments={props.apartments}
        favorites={props.favorites}
        select={props.select}
        toggleFavorite={props.toggleFavorite}
      />
    </section>
  );
}

function ApartmentGrid(props: {
  title: string;
  subtitle: string;
  apartments: Apartment[];
  favorites: Apartment[];
  select: (apartment: Apartment) => void;
  toggleFavorite: (apartment: Apartment) => void;
}) {
  return (
    <section className="results">
      <div className="sectionHead">
        <div>
          <h2>{props.title}</h2>
          <p>{props.subtitle}</p>
        </div>
      </div>
      <div className="cards">
        {props.apartments.length === 0 ? (
          <div className="emptyState">Пока ничего не найдено. Попробуй убрать часть фильтров.</div>
        ) : props.apartments.map((apartment) => (
          <article className="apartmentCard" key={apartment.id}>
            <button className="imageButton" onClick={() => props.select(apartment)}>
              <img src={apartment.photos[0] || fallbackApartments[0].photos[0]} alt={apartment.title} />
            </button>
            <div className="cardBody">
              <div className="priceRow">
                <strong>{formatPrice(apartment.price)}</strong>
                <button className="heart" onClick={() => props.toggleFavorite(apartment)}>
                  {props.favorites.some((item) => item.id === apartment.id) ? "В избранном" : "В избранное"}
                </button>
              </div>
              <h3>{apartment.title}</h3>
              <p>{apartment.district} р-н, {apartment.address}</p>
              <div className="chips">
                <span>{apartment.rooms} комн.</span>
                <span>{apartment.floor} этаж</span>
                {apartment.hasFurniture && <span>мебель</span>}
                {apartment.hasWifi && <span>Wi-Fi</span>}
              </div>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

function ApartmentDetails({ apartment, isFavorite, close, toggleFavorite, openConversation, deleteApartment, user }: {
  apartment: Apartment;
  isFavorite: boolean;
  close: () => void;
  toggleFavorite: (apartment: Apartment) => void;
  openConversation: (apartment: Apartment) => void;
  deleteApartment: (apartment: Apartment) => void;
  user: User | null;
}) {
  const canManage = user?.role === "landlord" && user.id === apartment.ownerId;

  return (
    <div className="modalLayer" onClick={close}>
      <section className="details" onClick={(event) => event.stopPropagation()}>
        <button className="close" onClick={close}>Закрыть</button>
        <img src={apartment.photos[0] || fallbackApartments[0].photos[0]} alt={apartment.title} />
        <div className="detailsBody">
          <p className="eyebrow">{apartment.district} · {apartment.type}</p>
          <h2>{apartment.title}</h2>
          <strong>{formatPrice(apartment.price)}</strong>
          <p>{apartment.description}</p>
          <div className="featureGrid">
            <span>{apartment.rooms} комнаты</span>
            <span>{apartment.floor} этаж</span>
            <span>{apartment.hasFurniture ? "С мебелью" : "Без мебели"}</span>
            <span>{apartment.hasWifi ? "Wi-Fi есть" : "Wi-Fi нет"}</span>
            <span>{apartment.hasWasher ? "Стиральная машина" : "Без стиральной машины"}</span>
          </div>
          <div className="actions">
            {!canManage && <button onClick={() => openConversation(apartment)}>Связаться</button>}
            <button className="secondary" onClick={() => toggleFavorite(apartment)}>
              {isFavorite ? "Убрать из избранного" : "Сохранить"}
            </button>
            {canManage && (
              <button className="danger" onClick={() => deleteApartment(apartment)}>
                Удалить объявление
              </button>
            )}
          </div>
        </div>
      </section>
    </div>
  );
}

function CreateApartmentForm({ user, onSubmit }: { user: User | null; onSubmit: (event: FormEvent<HTMLFormElement>) => void }) {
  if (!user) return <AuthRequired text="Войди как арендодатель, чтобы подать объявление." />;
  if (user.role !== "landlord") return <AuthRequired text="Публикация объявлений доступна только арендодателям." />;

  return (
    <section className="panel">
      <h2>Подать объявление</h2>
      <p>Заполни ключевые поля. Фото можно вставить ссылками, по одной строке.</p>
      <form className="formGrid" onSubmit={onSubmit}>
        <input name="title" placeholder="Заголовок" required />
        <select name="type" defaultValue="apartment">
          <option value="apartment">Квартира</option>
          <option value="studio">Студия</option>
          <option value="room">Комната</option>
          <option value="house">Дом</option>
        </select>
        <input name="price" placeholder="Цена" type="number" required />
        <select name="district" required>
          {districts.map((district) => <option key={district} value={district}>{district}</option>)}
        </select>
        <input name="address" placeholder="Адрес" required />
        <input name="rooms" placeholder="Комнаты" type="number" required />
        <input name="floor" placeholder="Этаж" type="number" required />
        <textarea name="description" placeholder="Описание" />
        <textarea name="photo_urls" placeholder="Ссылки на фото, каждая с новой строки" />
        <div className="checkStack inline">
          <label><input name="has_furniture" type="checkbox" /> Мебель</label>
          <label><input name="has_wifi" type="checkbox" /> Wi-Fi</label>
          <label><input name="has_washer" type="checkbox" /> Стиральная машина</label>
        </div>
        <button className="wide">Опубликовать</button>
      </form>
    </section>
  );
}

function Chats(props: {
  conversations: Conversation[];
  activeConversation: number | null;
  messages: Message[];
  openConversation: (id: number) => void;
  sendMessage: (event: FormEvent<HTMLFormElement>) => void;
  user: User | null;
}) {
  if (!props.user) return <AuthRequired text="Войди, чтобы видеть чаты." />;

  return (
    <section className="chatLayout">
      <aside className="chatList">
        <h2>Мои чаты</h2>
        {props.conversations.length === 0 ? (
          <p>Диалогов пока нет. Открой карточку и нажми “Связаться”.</p>
        ) : props.conversations.map((conversation) => (
          <button
            key={conversation.id}
            className={props.activeConversation === conversation.id ? "active" : ""}
            onClick={() => props.openConversation(conversation.id)}
          >
            <strong>{conversation.apartment?.title || `Диалог #${conversation.id}`}</strong>
            <span>{conversation.apartment?.district || "AlaRent"}</span>
          </button>
        ))}
      </aside>
      <section className="messages">
        <div className="messagesBody">
          {props.messages.length === 0 ? (
            <p className="emptyState">Выбери диалог или напиши первое сообщение.</p>
          ) : props.messages.map((message) => (
            <div key={message.id} className={message.senderId === props.user?.id ? "message own" : "message"}>
              {message.text}
            </div>
          ))}
        </div>
        <form className="messageForm" onSubmit={props.sendMessage}>
          <input name="text" placeholder="Написать сообщение" />
          <button>Отправить</button>
        </form>
      </section>
    </section>
  );
}

function Profile({ user, authMode, setAuthMode, handleAuth, updateProfile, notice }: {
  user: User | null;
  authMode: "login" | "register";
  setAuthMode: (mode: "login" | "register") => void;
  handleAuth: (event: FormEvent<HTMLFormElement>) => void;
  updateProfile: (event: FormEvent<HTMLFormElement>) => void;
  notice: string;
}) {
  if (!user) {
    return (
      <section className="authPage">
        <div className="authVisual">
          <p className="eyebrow">AlaRent аккаунт</p>
          <h1>{authMode === "login" ? "Вернись к сохраненным вариантам." : "Создай профиль для аренды."}</h1>
          <p>
            Один аккаунт для поиска жилья, публикации объявлений, избранного и быстрых
            диалогов с арендодателями.
          </p>
          <div className="authBenefits">
            <span>Избранное</span>
            <span>Чаты</span>
            <span>Объявления</span>
          </div>
        </div>

        <div className="authCard">
          <div className="authTabs">
            <button className={authMode === "login" ? "active" : ""} onClick={() => setAuthMode("login")}>
              Вход
            </button>
            <button className={authMode === "register" ? "active" : ""} onClick={() => setAuthMode("register")}>
              Регистрация
            </button>
          </div>

          <div className="authCardHead">
            <h2>{authMode === "login" ? "Войти в аккаунт" : "Новый аккаунт"}</h2>
            <p>{authMode === "login" ? "Продолжи поиск с того места, где остановился." : "Выбери роль и заполни базовые данные."}</p>
          </div>

          {notice && <div className="authNotice">{notice}</div>}

          <form className="authForm" onSubmit={handleAuth}>
            {authMode === "register" && (
              <>
                <input name="name" placeholder="Имя" required />
                <input name="phone" placeholder="Телефон" required />
                <select name="role" defaultValue="tenant">
                  <option value="tenant">Ищу жилье</option>
                  <option value="landlord">Сдаю жилье</option>
                </select>
              </>
            )}
            <input name="email" placeholder="Email" type="email" required />
            <input name="password" placeholder="Пароль" type="password" required />
            <button>{authMode === "login" ? "Войти" : "Создать аккаунт"}</button>
          </form>

          <button className="authSwitch" onClick={() => setAuthMode(authMode === "login" ? "register" : "login")}>
            {authMode === "login" ? "Нет аккаунта? Зарегистрироваться" : "Уже есть аккаунт? Войти"}
          </button>
        </div>
      </section>
    );
  }

  return (
    <section className="profilePage">
      <div className="profileHero">
        <div className="avatar">{user.name.slice(0, 1).toUpperCase()}</div>
        <div>
          <p className="eyebrow">{user.role === "landlord" ? "Арендодатель" : "Арендатор"}</p>
          <h1>{user.name}</h1>
          <p>{user.email}</p>
        </div>
      </div>
      <form className="profileForm" onSubmit={updateProfile}>
        <input name="name" defaultValue={user.name} placeholder="Имя" />
        <input name="phone" defaultValue={user.phone} placeholder="Телефон" />
        <input name="avatar_url" defaultValue={user.avatarUrl || ""} placeholder="Фото профиля URL" />
        <textarea name="bio" defaultValue={user.bio || ""} placeholder="О себе" />
        <button className="wide">Сохранить профиль</button>
      </form>
    </section>
  );
}

function AuthRequired({ text }: { text: string }) {
  return (
    <section className="panel emptyState">
      <h2>Нужен доступ</h2>
      <p>{text}</p>
    </section>
  );
}

function buildFilterParams(filters: Filters) {
  const params = new URLSearchParams();
  if (filters.q) params.set("q", filters.q);
  if (filters.district) params.set("district", filters.district);
  if (filters.minPrice) params.set("min_price", filters.minPrice);
  if (filters.maxPrice) params.set("max_price", filters.maxPrice);
  if (filters.rooms) params.set("rooms", filters.rooms);
  if (filters.hasWifi) params.set("has_wifi", "true");
  if (filters.hasFurniture) params.set("has_furniture", "true");
  if (filters.hasWasher) params.set("has_washer", "true");
  if (filters.sort) params.set("sort", filters.sort);
  const value = params.toString();
  return value ? `?${value}` : "";
}

function normalizeApartment(raw: unknown): Apartment {
  const item = raw as Record<string, unknown>;
  const photosRaw = (read(item, "Photos", "photos") as unknown[]) || [];
  const photos = Array.isArray(photosRaw)
    ? photosRaw.map((photo) => String(read(photo, "URL", "url") || photo)).filter(Boolean)
    : [];

  return {
    id: Number(read(item, "ID", "id") || 0),
    title: String(read(item, "Title", "title") || "Без названия"),
    description: String(read(item, "Description", "description") || ""),
    type: String(read(item, "Type", "type") || "apartment"),
    price: Number(read(item, "Price", "price") || 0),
    district: String(read(item, "District", "district") || "Астана"),
    address: String(read(item, "Address", "address") || ""),
    rooms: Number(read(item, "Rooms", "rooms") || 1),
    floor: Number(read(item, "Floor", "floor") || 1),
    hasFurniture: Boolean(read(item, "HasFurniture", "has_furniture", "hasFurniture")),
    hasWifi: Boolean(read(item, "HasWifi", "has_wifi", "hasWifi")),
    hasWasher: Boolean(read(item, "HasWasher", "has_washer", "hasWasher")),
    ownerId: Number(read(item, "OwnerID", "owner_id", "ownerId") || 0),
    photos,
  };
}

function normalizeUser(raw: unknown): User {
  const item = raw as Record<string, unknown>;
  return {
    id: Number(read(item, "ID", "id") || 0),
    name: String(read(item, "Name", "name") || "Пользователь"),
    email: String(read(item, "Email", "email") || ""),
    phone: String(read(item, "Phone", "phone") || ""),
    role: String(read(item, "Role", "role") || "tenant") as Role,
    avatarUrl: String(read(item, "AvatarURL", "avatar_url", "avatarUrl") || ""),
    bio: String(read(item, "Bio", "bio") || ""),
  };
}

function normalizeConversation(raw: unknown): Conversation {
  const item = raw as Record<string, unknown>;
  const apartment = read(item, "Apartment", "apartment");
  return {
    id: Number(read(item, "ID", "id") || 0),
    apartment: apartment ? normalizeApartment(apartment) : undefined,
    tenant: read(item, "Tenant", "tenant") ? normalizeUser(read(item, "Tenant", "tenant")) : undefined,
    landlord: read(item, "Landlord", "landlord") ? normalizeUser(read(item, "Landlord", "landlord")) : undefined,
  };
}

function normalizeMessage(raw: unknown): Message {
  const item = raw as Record<string, unknown>;
  return {
    id: Number(read(item, "ID", "id") || Math.random()),
    text: String(read(item, "Text", "text") || ""),
    senderId: Number(read(item, "SenderID", "sender_id", "senderId") || 0),
    createdAt: String(read(item, "CreatedAt", "created_at", "createdAt") || ""),
  };
}

function read(raw: unknown, ...keys: string[]) {
  const item = raw as Record<string, unknown>;
  return keys.find((key) => item && item[key] !== undefined) ? item[keys.find((key) => item[key] !== undefined)!] : undefined;
}

function readStoredUser() {
  try {
    const raw = localStorage.getItem("alarent_user");
    return raw ? normalizeUser(JSON.parse(raw)) : null;
  } catch {
    return null;
  }
}

function formatPrice(value: number) {
  return new Intl.NumberFormat("ru-KZ").format(value) + " ₸ / мес";
}

function getErrorMessage(error: unknown) {
  return error instanceof Error ? error.message : "Что-то пошло не так";
}

export default App;
