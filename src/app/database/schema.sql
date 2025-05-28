--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Debian 15.13-1.pgdg120+1)
-- Dumped by pg_dump version 15.13 (Debian 15.13-0+deb12u1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: estado_producto; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.estado_producto AS ENUM (
    'disponible',
    'agotado',
    'discontinuado'
);


ALTER TYPE public.estado_producto OWNER TO postgres;

--
-- Name: nombre_categoria; Type: DOMAIN; Schema: public; Owner: postgres
--

CREATE DOMAIN public.nombre_categoria AS character varying(50)
	CONSTRAINT nombre_categoria_check CHECK (((VALUE)::text ~* '^[a-zA-Z0-9 \u00C0-\u00FF]+$'::text));


ALTER DOMAIN public.nombre_categoria OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categorias; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categorias (
    id bigint NOT NULL,
    nombre public.nombre_categoria NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.categorias OWNER TO postgres;

--
-- Name: categorias_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categorias_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.categorias_id_seq OWNER TO postgres;

--
-- Name: categorias_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categorias_id_seq OWNED BY public.categorias.id;


--
-- Name: productos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.productos (
    id bigint NOT NULL,
    nombre character varying(100) NOT NULL,
    precio numeric(10,2) NOT NULL,
    descripcion text,
    estado public.estado_producto DEFAULT 'disponible'::public.estado_producto NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.productos OWNER TO postgres;

--
-- Name: productos_categorias; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.productos_categorias (
    producto_id bigint NOT NULL,
    categoria_id bigint NOT NULL
);


ALTER TABLE public.productos_categorias OWNER TO postgres;

--
-- Name: productos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.productos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.productos_id_seq OWNER TO postgres;

--
-- Name: productos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.productos_id_seq OWNED BY public.productos.id;


--
-- Name: vista_producto_categoria; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.vista_producto_categoria AS
 SELECT p.id AS producto_id,
    p.nombre AS producto_nombre,
    p.precio,
    p.descripcion AS producto_descripcion,
    p.estado AS producto_estado,
    c.id AS categoria_id,
    c.nombre AS categoria_nombre
   FROM ((public.productos p
     JOIN public.productos_categorias pc ON ((p.id = pc.producto_id)))
     JOIN public.categorias c ON ((c.id = pc.categoria_id)))
  WHERE ((p.deleted_at IS NULL) AND (c.deleted_at IS NULL));


ALTER TABLE public.vista_producto_categoria OWNER TO postgres;

--
-- Name: categorias id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categorias ALTER COLUMN id SET DEFAULT nextval('public.categorias_id_seq'::regclass);


--
-- Name: productos id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.productos ALTER COLUMN id SET DEFAULT nextval('public.productos_id_seq'::regclass);


--
-- Data for Name: categorias; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categorias (id, nombre, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: productos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.productos (id, nombre, precio, descripcion, estado, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: productos_categorias; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.productos_categorias (producto_id, categoria_id) FROM stdin;
\.


--
-- Name: categorias_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categorias_id_seq', 1, false);


--
-- Name: productos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.productos_id_seq', 1, false);


--
-- Name: categorias categorias_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categorias
    ADD CONSTRAINT categorias_pkey PRIMARY KEY (id);


--
-- Name: productos_categorias productos_categorias_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.productos_categorias
    ADD CONSTRAINT productos_categorias_pkey PRIMARY KEY (producto_id, categoria_id);


--
-- Name: productos productos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.productos
    ADD CONSTRAINT productos_pkey PRIMARY KEY (id);


--
-- Name: categorias uni_categorias_nombre; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categorias
    ADD CONSTRAINT uni_categorias_nombre UNIQUE (nombre);


--
-- Name: idx_categorias_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categorias_deleted_at ON public.categorias USING btree (deleted_at);


--
-- Name: idx_productos_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_productos_deleted_at ON public.productos USING btree (deleted_at);


--
-- Name: productos_categorias fk_productos_categorias_categoria; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.productos_categorias
    ADD CONSTRAINT fk_productos_categorias_categoria FOREIGN KEY (categoria_id) REFERENCES public.categorias(id);


--
-- Name: productos_categorias fk_productos_categorias_producto; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.productos_categorias
    ADD CONSTRAINT fk_productos_categorias_producto FOREIGN KEY (producto_id) REFERENCES public.productos(id);


--
-- PostgreSQL database dump complete
--

